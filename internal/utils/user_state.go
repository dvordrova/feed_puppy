package utils

type UserState string

const (
	UserStateNewDog      UserState = "new_dog"
	UserStateActionFeed  UserState = "action_feed"
	UserStateActionYum   UserState = "action_yum"
	UserStateActionWeigh UserState = "action_weigh"
	UserStateDogSelected UserState = "dog_selected"
	UserStateChangeName  UserState = "change_name"
)
