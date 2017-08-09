package store

import (
	"testing"
	"testing/quick"

	"github.com/trussle/snowy/pkg/uuid"
)

func TestVirtualStore(t *testing.T) {
	t.Parallel()

	t.Run("get", func(t *testing.T) {
		store := NewVirtualStore()
		_, err := store.Get(uuid.New(), Query{})

		if expected, actual := true, err != nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		if expected, actual := true, ErrNotFound(err); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("get when empty", func(t *testing.T) {
		var (
			id    = uuid.New()
			store = NewVirtualStore()
		)

		if err := store.Insert(Entity{ResourceID: id}); err != nil {
			t.Fatal(err)
		}

		if v, ok := store.(*virtualStore); ok {
			v.entities[id.String()] = make([]Entity, 0)
		}

		_, err := store.Get(id, Query{})

		if expected, actual := true, err != nil; expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}

		if expected, actual := true, ErrNotFound(err); expected != actual {
			t.Errorf("expected: %t, actual: %t", expected, actual)
		}
	})

	t.Run("put", func(t *testing.T) {
		store := NewVirtualStore()

		fn := func(res uuid.UUID) bool {
			err := store.Insert(Entity{ResourceID: res})
			return err == nil
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put then get", func(t *testing.T) {
		store := NewVirtualStore()

		fn := func(res uuid.UUID) bool {
			if err := store.Insert(Entity{ResourceID: res}); err != nil {
				t.Fatal(err)
			}

			entity, err := store.Get(res, Query{})
			if err != nil {
				t.Fatal(err)
			}

			return entity.ResourceID == res
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put then query with no id", func(t *testing.T) {
		store := NewVirtualStore()

		fn := func(res uuid.UUID) bool {
			if err := store.Insert(Entity{ResourceID: res}); err != nil {
				t.Fatal(err)
			}

			entities, err := store.GetMultiple(uuid.New(), Query{
				Tags: make([]string, 0),
			})
			if err != nil {
				t.Fatal(err)
			}

			return len(entities) == 0
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put then query with no tags should select all", func(t *testing.T) {
		store := NewVirtualStore()

		fn := func(res uuid.UUID) bool {
			if err := store.Insert(Entity{ResourceID: res}); err != nil {
				t.Fatal(err)
			}

			entities, err := store.GetMultiple(res, Query{
				Tags: make([]string, 0),
			})
			if err != nil {
				t.Fatal(err)
			}

			return len(entities) == 1
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put then query exact match", func(t *testing.T) {
		store := NewVirtualStore()

		fn := func(res uuid.UUID, tags Tags) bool {
			entity := Entity{
				ResourceID: res,
				Tags:       tags.Slice(),
			}
			if err := store.Insert(entity); err != nil {
				t.Fatal(err)
			}

			got, err := store.GetMultiple(res, Query{
				Tags: tags,
			})
			if err != nil {
				t.Fatal(err)
			}

			want := []Entity{entity}
			return equals(want, got)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("multiple puts then query exact match", func(t *testing.T) {
		store := NewVirtualStore()

		fn := func(res uuid.UUID, tags Tags) bool {
			want := make([]Entity, 10)
			for k := range want {
				entity := Entity{
					ResourceID: res,
					Tags:       tags.Slice(),
				}
				if err := store.Insert(entity); err != nil {
					t.Fatal(err)
				}
				want[k] = entity
			}

			got, err := store.GetMultiple(res, Query{
				Tags: tags,
			})
			if err != nil {
				t.Fatal(err)
			}

			return equals(want, got)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("puts then query partial match", func(t *testing.T) {
		store := NewVirtualStore()

		fn := func(res uuid.UUID, tags Tags) bool {
			want := make([]Entity, 10)
			for k := range want {
				entity := Entity{
					ResourceID: res,
					Tags:       tags.Slice(),
				}
				if err := store.Insert(entity); err != nil {
					t.Fatal(err)
				}
				want[k] = entity
			}

			got, err := store.GetMultiple(res, Query{
				Tags: splitTags(tags.Slice()),
			})
			if err != nil {
				t.Fatal(err)
			}

			return equals(want, got)
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("put then query with no tags should not select any", func(t *testing.T) {
		store := NewVirtualStore()

		fn := func(res uuid.UUID) bool {
			if err := store.Insert(Entity{
				ResourceID: res,
				Tags:       []string{"a"},
			}); err != nil {
				t.Fatal(err)
			}

			entities, err := store.GetMultiple(res, Query{
				Tags: []string{"b"},
			})
			if err != nil {
				t.Fatal(err)
			}

			return len(entities) == 0
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("drop", func(t *testing.T) {
		store := NewVirtualStore()

		fn := func(res uuid.UUID) bool {
			if err := store.Insert(Entity{ResourceID: res}); err != nil {
				t.Fatal(err)
			}

			entity, err := store.Get(res, Query{})
			if err != nil {
				t.Fatal(err)
			}

			if expected, actual := res, entity.ResourceID; !expected.Equals(actual) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}

			if err = store.Drop(); err != nil {
				t.Fatal(err)
			}

			if expected, actual := false, ErrNotFound(err); expected != actual {
				t.Errorf("expected: %t, actual: %t", expected, actual)
			}

			return true
		}

		if err := quick.Check(fn, nil); err != nil {
			t.Error(err)
		}
	})

	t.Run("run and stop", func(t *testing.T) {
		store := NewVirtualStore()

		if err := store.Run(); err != nil {
			t.Error(err)
		}

		store.Stop()
	})
}
