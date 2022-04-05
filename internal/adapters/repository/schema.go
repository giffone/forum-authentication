package repository

import (
	"forum/internal/constant"
)

// Configuration for save parameters
type Configuration struct {
	Name, Path, PathB,
	Driver, Port, Connection string
}

func MakeTables() []string {
	return []string{
		constant.TabUsers,
		constant.TabCategories,
		constant.TabLikes,
		constant.TabPosts,
		constant.TabPostsLikes,
		constant.TabPostsCategories,
		constant.TabComments,
		constant.TabCommentsLikes,
		constant.TabSessions}
}


