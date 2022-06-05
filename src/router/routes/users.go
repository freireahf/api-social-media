package routes

import (
	"api/src/controllers"
	"net/http"
)

var userRoutes = []Route{
	{
		URI:                   "/users",
		Method:                http.MethodPost,
		Function:              controllers.CreateUser,
		RequireAuthentication: false,
	},
	{
		URI:                   "/users",
		Method:                http.MethodGet,
		Function:              controllers.FindAllUsersFilteredByNameOrNick,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}",
		Method:                http.MethodGet,
		Function:              controllers.FindUserById,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}",
		Method:                http.MethodPut,
		Function:              controllers.UpdateUserById,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}",
		Method:                http.MethodDelete,
		Function:              controllers.DeleteUser,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}/follower",
		Method:                http.MethodPost,
		Function:              controllers.FollowerUser,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}/unfollow",
		Method:                http.MethodPost,
		Function:              controllers.UnfollowUser,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}/followers",
		Method:                http.MethodGet,
		Function:              controllers.FindFollowers,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}/following",
		Method:                http.MethodGet,
		Function:              controllers.FindFollowing,
		RequireAuthentication: true,
	},
	{
		URI:                   "/users/{userId}/password",
		Method:                http.MethodPost,
		Function:              controllers.UpdatePassword,
		RequireAuthentication: true,
	},
}
