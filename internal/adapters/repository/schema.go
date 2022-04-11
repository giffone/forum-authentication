package repository

import (
	"github.com/giffone/forum-authentication/internal/constant"
)

// Configuration for save parameters
type Configuration struct {
	Name, Path, PathB,
	Driver, Port, Connection string
}

// MakeTables returns a list of tables that need to create
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
		constant.TabSessions,
		constant.TabTokens,
	}
}
