package constant

import (
	"time"
)

const (
	URLHome         = "/"               // homepage
	URLSignUp       = "/signup"         // for create user (in login page)
	URLLogin        = "/login"          // for begin session
	URLLogout       = "/logout"         // for end session
	URLPost         = "/post"           // create post
	URLRead         = "/read/"          // view one post
	URLReadRatio    = "/read/ratio/"    // create like
	URLCategory     = "/category"       // create category
	URLCategoryBy   = "/category/"      // homepage sorted by categories
	URLComment      = "/comment"        // create comment
	URLAccount      = "/account/"       // administrator page
	URLAccountRatio = "/account/ratio/" // administrator page
	URLFavicon      = "/favicon.ico"    // add favicon site

	/*------------------------------------------------------*/

	PathIndex         = "internal/web/templates/index.gohtml"
	PathPost          = "internal/web/templates/post.gohtml"
	PathLoginObj      = "internal/web/templates/login.gohtml"
	PathAccount       = "internal/web/templates/account.gohtml"
	PathMessage       = "internal/web/templates/message.gohtml"
	PathAccountUser   = "internal/web/templates/object/account.gohtml"
	PathPostObj       = "internal/web/templates/object/post.gohtml"
	PathHeaderObj     = "internal/web/templates/object/header.gohtml"
	PathCategoriesObj = "internal/web/templates/object/categories.gohtml"
	PathLikesObj      = "internal/web/templates/object/likes.gohtml"
	PathCommentsObj   = "internal/web/templates/object/comments.gohtml"
	PathFooterObj     = "internal/web/templates/object/footer.gohtml"

	/*------------------------------------------------------*/

	CookieSession     = "session"         // name for cookie
	CookieUserID      = "userID"          // name for cookie
	CookiePostID      = "postID"          // name for cookie
	LoginMinLength    = 2                 // symbols
	PasswordMinLength = 6                 // symbols
	PostShowOnPage    = 10                // 10 post will show on main page
	SessionExpire     = 1                 // 1 day (in days)
	SessionMaxAge     = 24 * 60 * 60      // 1 day (in seconds)
	TimeLimit         = 10 * time.Second  // context (for handlers, including queries to database)
	TimeLimitDB       = 5 * time.Second   // context (for queries to database)
	ForumLayoutDate   = "January 2, 2006" // format for page

	/*------------------------------------------------------*/

	Code200 = 200 // http.StatusOK (GET)
	Code201 = 201 // http.StatusCreated (POST)
	Code204 = 204 // http.StatusNoContent (PUT,PATCH,DELETE)
	Code302 = 302 // http.StatusFound
	Code400 = 400 // http.StatusBadRequest
	Code401 = 401 // http.StatusUnauthorized
	Code403 = 403 // http.StatusForbidden
	Code404 = 404 // http.StatusNotFound
	Code405 = 405 // http.StatusMethodNotAllowed
	Code422 = 422 // http.StatusUnprocessableEntity
	Code500 = 500 // http.StatusInternalServerError

	/*------------------------------------------------------*/

	StatusOK          = "Successfully: %s"
	StatusCreated     = "Created: %s"
	AlreadyExist      = "can not create: %s already have"
	InvalidCharacters = "invalid: the %s contains invalid characters"
	TooShort          = "too short:  must be at least %s characters"
	NotMatch          = "no match: the entered %s does not match"
	WrongEnter        = "wrong: the entered %s is wrong"
	InvalidEnter      = "invalid: the entered %s is incorrect, please use valid"
	InternalError     = "internal error: \"%v\""
	AccessDenied      = "access denied: you not authorized or session expired"

	/*------------------------------------------------------*/

	TabUsers                    = "src_users"
	TabCategories               = "src_categories"
	TabLikes                    = "src_likes"
	TabPosts                    = "posts"
	TabPostsLikes               = "posts_likes"
	TabPostsCategories          = "posts_categories"
	TabComments                 = "comments"
	TabCommentsLikes            = "comments_likes"
	TabSessions                 = "sessions"
	QueAttach                   = "attach"
	QueDetach                   = "detach"
	QueRestore                  = "restore"
	QueInsert2                  = "insert_2"
	QueInsert3                  = "insert_3"
	QueInsert4                  = "insert_4"
	QueInsert5                  = "insert_5"
	QueSelect                   = "select"       // all posts without any sort
	QueSelectPosts              = "select_posts" // all posts without any sort
	QueSelectUsers              = "select_users"
	QueSelectPostsBy            = "select_posts_by" // all posts sorted by WHERE
	QueSelectPostsRatedBy       = "select_posts_rated_by"
	QueSelectCommentsRatedBy    = "select_comments_rated_by"
	QueSelectPostsAndCategoryBy = "select_posts_category_by"
	QueSelectCommentAndPostBy   = "select_posts_comment_by"
	QueSelectCategories         = "select_categories"
	QueSelectUserBy             = "select_user_by"
	QueSelectCategoryBy         = "select_category_by"
	QueSelectSessionBy          = "select_session_by"
	QueSelectLikeCountBy        = "select_post_like_count_by"
	QueSelectCommentLikeCountBy = "select_comment_like_count_by"
	QueSelectCommentsBy         = "select_comments_by"
	QueSelectLikeBy             = "select_like_by"
	QueDeleteBy                 = "delete_session_by"
	QueSelectCount              = "select_count"
	QueSelectLikedOrNot         = "select_liked_or_not"

	/*------------------------------------------------------*/

	FieldID           = "id"
	FieldLike         = "like"
	FieldUUID         = "uuid"
	FieldUser         = "user"
	FieldPost         = "post"
	FieldBody         = "body"
	FieldRoot         = "root"
	FieldLiked        = "liked"
	FieldLikes        = "likes"
	FieldLogin        = "login"
	FieldTitle        = "title"
	FieldEmail        = "email"
	FieldExpire       = "expire"
	FieldDislike      = "dislike"
	FieldCreated      = "created"
	FieldUUIDID       = "uuid-id"
	FieldComment      = "comment"
	FieldCategory     = "category"
	FieldCategories   = "categories"
	FieldPassword     = "password"
	FieldPostLiked    = "post-liked"
	FieldCommentLiked = "comment-liked"

	/*------------------------------------------------------*/

	KeyID           = "id"
	KeyLogin        = "login"
	KeyUser         = "user"
	KeyPost         = "post"
	KeyComment      = "comment"
	KeyPostRated    = "post rated"
	KeyCommentRated = "comment rated"
	KeyLike         = "like"
	KeyLink         = "link"
	KeyDislike      = "dislike"
	KeyRated        = "rated"
	KeyCategory     = "category"
	KeyRate         = "rate"
	KeyObject       = "object"
)
