package git

import (
	"sync"
)

func LoadRepositoryEntities(directories []string) (entities []*RepoEntity, err error) {
	entities = make([]*RepoEntity, 0)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, dir := range directories {
		// increment wait counter by one because we run a single goroutine
		// below
		wg.Add(1)
		go func(d string) {
			// decrement the wait counter by one, we call it in a defer so it's
			// called at the end of this goroutine
			defer wg.Done()
			entity, err := InitializeRepository(d)
			if err != nil {
				return
			}
			// lock so we don't get a race if multiple go routines try to add
			// to the same entities
			mu.Lock()
			entities = append(entities, entity)
			mu.Unlock()
		}(dir)
	}
	// wait until the wait counter is zero, this happens if all goroutines have
	// finished
	wg.Wait()
	return entities, nil
}