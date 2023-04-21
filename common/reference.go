package common

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
	"google.golang.org/protobuf/proto"
)

var References *ReferenceManager

type ReferenceItemManager[T proto.Message] struct {
	manager      *ReferenceManager
	refType      ReferenceType
	path         string
	cache        cmap.ConcurrentMap[string, T]
	loadingAll   *sync.WaitGroup
	newItemFn    func() T
	unmarshallFn func(data []byte) (T, error)
}

type ReferenceClient struct {
}

type ReferenceManager struct {
	loadingAll *sync.WaitGroup

	Weapons    *ReferenceItemManager[*RefWeapon]
	GameTerms  *ReferenceItemManager[*RefGameTerm]
	Skills     *ReferenceItemManager[*RefSkill]
	Characters *ReferenceItemManager[*Character]
	Armors     *ReferenceItemManager[*RefArmor]
}

func newReferenceManager() *ReferenceManager {
	r := &ReferenceManager{}
	r.loadingAll = new(sync.WaitGroup)

	r.Weapons = newReferenceItemManager(ReferenceType_ReferenceType_Weapon, "weapons", func(data []byte) (*RefWeapon, error) {
		w4 := &RefWeapon{}
		err := proto.Unmarshal(data, w4)
		return w4, err
	})

	r.GameTerms = newReferenceItemManager(ReferenceType_ReferenceType_GameTerm, "gameterms", func(data []byte) (*RefGameTerm, error) {
		item := &RefGameTerm{}
		err := proto.Unmarshal(data, item)
		return item, err
	})

	r.Skills = newReferenceItemManager(ReferenceType_ReferenceType_Skill, "skills", func(data []byte) (*RefSkill, error) {
		item := &RefSkill{}
		err := proto.Unmarshal(data, item)
		return item, err
	})

	r.Characters = newReferenceItemManager(ReferenceType_ReferenceType_Skill, "characters", func(data []byte) (*Character, error) {
		item := &Character{}
		err := proto.Unmarshal(data, item)
		return item, err
	})

	r.Armors = newReferenceItemManager(ReferenceType_ReferenceType_Skill, "armors", func(data []byte) (*RefArmor, error) {
		item := &RefArmor{}
		err := proto.Unmarshal(data, item)
		return item, err
	})

	return r
}

func (r *ReferenceManager) LoadAll() {
	r.loadingAll.Add(3)
	r.Weapons.LoadAll(r.loadingAll)
	r.Skills.LoadAll(r.loadingAll)
	r.GameTerms.LoadAll(r.loadingAll)

}

func newReferenceItemManager[T proto.Message](refType ReferenceType, path string, unmarshall func(data []byte) (T, error)) *ReferenceItemManager[T] {
	r := &ReferenceItemManager[T]{
		path:         path,
		refType:      refType,
		unmarshallFn: unmarshall,
	}
	r.cache = cmap.New[T]()
	r.loadingAll = new(sync.WaitGroup)
	return r
}

func (r *ReferenceItemManager[T]) ExtractID(path string) string {
	p := filepath.Base(path)
	if strings.HasSuffix(p, ".pbf") {
		return p[:len(p)-4]
	}
	return p
}

func (r *ReferenceItemManager[T]) PathTo(ID string) string {
	if ID == "" {
		return filepath.Join("reference", r.path)
	}
	if strings.HasPrefix(ID, ".pbf") {
		return filepath.Join("reference", r.path, ID)
	}
	return filepath.Join("reference", r.path, fmt.Sprintf("%v.pbf", ID))
}

func (r *ReferenceItemManager[T]) LoadAll(wg *sync.WaitGroup) {
	r.loadingAll.Add(1)
	path := r.PathTo("")
	keys, err := Storage.GetAllKeys(path)
	if err != nil {
		log.Fatalf("Error loading reference items %v, %v", r.path, err)
	}
	for _, k := range keys {
		id := r.ExtractID(k)
		r.Load(id)
	}
	r.loadingAll.Done()
	wg.Done()
}

func (r *ReferenceItemManager[T]) Load(id string) T {
	key := r.PathTo(id)
	data, err := Storage.Get(key)
	if err != nil {
		log.Fatalf("Error loading reference items %v, %v", r.path, err)
	}

	item, err := r.unmarshallFn(data)
	if err != nil {
		log.Fatalf("Error unmarshalling reference items %v, %v", r.path, err)
	}

	r.cache.Set(id, item)
	// var w Weapon
	return item
}

func (r *ReferenceItemManager[T]) All() []T {
	rtn := make([]T, r.cache.Count())
	cnt := 0
	for item := range r.cache.IterBuffered() {
		rtn[cnt] = item.Val
		cnt++
	}
	return rtn
}

func (r *ReferenceItemManager[T]) Count() int {
	return r.cache.Count()
}

func (r *ReferenceItemManager[T]) Get(id string) T {
	item, found := r.cache.Get(id)
	if found {
		return item
	}
	return r.Load(id)
}

func (r *ReferenceItemManager[T]) Set(item T, id string) {
	data, err := proto.Marshal(item)
	if err != nil {
		log.Fatalf("Error Marshalling: %v", err)
	}
	key := r.PathTo(id)
	err = Storage.Put(key, data)
	if err != nil {
		log.Fatalf("Error Saving: %v", err)
	}
	r.cache.Set(id, item)
}

func (r *ReferenceItemManager[T]) Delete(id string) {
	path := r.PathTo(id)
	err := Storage.Delete(path)
	if err != nil {
		log.Fatalf("Error Deleteing: %v", err)
	}
}

// Client Side
func (r *ReferenceClient) GetWeapon(ctx context.Context, id string) *RefWeapon {
	request := &ReferenceRequest{
		Type: ReferenceType_ReferenceType_Weapon,
		ID:   id,
	}
	m := SendAndWait(ctx, &Fp2Message{
		Data: &Fp2Message_ReferenceRequest{
			ReferenceRequest: request,
		},
	})

	response := m.GetReferenceResponse()
	return response.GetWeapon()
}

// Server
func (r *ReferenceManager) Handle(m *Fp2Message) {
	refRequest := m.GetReferenceRequest()
	if refRequest != nil {
		switch refRequest.Type {

		// Weapons
		case ReferenceType_ReferenceType_Weapon:
			w := r.GetWeapon(refRequest.ID)
			resp := &ReferenceResponse{
				Reference: &ReferenceResponse_Weapon{
					Weapon: w,
				},
			}
			r.respond(resp, m.MessageID)

		// Game Terms
		case ReferenceType_ReferenceType_GameTerm:

		default:
			// OH NO
		}

	}
}

func (r *ReferenceManager) respond(resp *ReferenceResponse, requestID string) {
	m := &Fp2Message{
		RespondingToID: requestID,
		Data: &Fp2Message_ReferenceResponse{
			ReferenceResponse: resp,
		},
	}
	prepMessage(m)
	Comms.Send(context.Background(), m)
}

func (r *ReferenceManager) GetWeapon(ID string) *RefWeapon {
	path := fmt.Sprintf("reference/weapons/reference-weapon-%v.pbf", ID)
	data, err := Storage.Get(path)
	if err != nil {
		log.Printf("Error in GetWeapon: %v, %v\n", ID, err)
		return nil
	}
	w := &RefWeapon{}
	err = proto.Unmarshal(data, w)
	if err != nil {
		log.Printf("Invalid Data in GetWeapon: %v, %v\n", ID, err)
		return nil
	}
	return w
}

func (r ReferenceManager) SaveWeapon(w *RefWeapon) error {
	key := fmt.Sprintf("reference/weapons/reference-weapon-%v.pbf", w.ID)
	data, err := proto.Marshal(w)
	if err != nil {
		return err
	}
	return Storage.Put(key, data)
}
