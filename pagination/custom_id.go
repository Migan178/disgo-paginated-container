package pagination

import "fmt"

const (
	paginationContainerFirst   = "/paginated/container/first"
	paginationContainerPrev    = "/paginated/container/prev"
	paginationContainerPages   = "/paginated/container/pages"
	paginationContainerNext    = "/paginated/container/next"
	paginationContainerLast    = "/paginated/container/last"
	paginationContainerModal   = "/paginated/container/modal"
	paginationContainerSetPage = "/paginated/container/modal/set"
)

func makePaginationContainerPrev(id string) string {
	return fmt.Sprintf("%s/%s", paginationContainerPrev, id)
}

func makePaginationContainerFirst(id string) string {
	return fmt.Sprintf("%s/%s", paginationContainerFirst, id)
}

func makePaginationContainerPages(id string) string {
	return fmt.Sprintf("%s/%s", paginationContainerPages, id)
}

func makePaginationContainerNext(id string) string {
	return fmt.Sprintf("%s/%s", paginationContainerNext, id)
}

func makePaginationContainerLast(id string) string {
	return fmt.Sprintf("%s/%s", paginationContainerLast, id)
}

func makePaginationContainerModal(id string) string {
	return fmt.Sprintf("%s/%s", paginationContainerModal, id)
}
