package model

type VideoVisibility string
type CommentStatus string
type ReactionType string
type NotificationType string
type NotificationEntityType string
const (
	Public VideoVisibility = "public"
	Private VideoVisibility = "private"
)
const (
	Visible CommentStatus = "visible"
	Hidden CommentStatus = "hidden"
	Deleted CommentStatus = "deleted"
)
const (
	Like ReactionType = "like"
	Dislike ReactionType = "dislike"
)
const(
	NewVideo NotificationType = "new_video"
	CommentReply NotificationType = "comment_reply"
	System NotificationType = "system"
)
const(
	VideoEntity NotificationEntityType = "video"
	CommentEntity NotificationEntityType = "comment"
)