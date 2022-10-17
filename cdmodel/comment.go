package cdmodel

type Comment struct {
	CommentID		int
	Content			string
	Time			string
	Author			string
	Self			bool
	Replies			[]CommentReply
	Score			int
}

type CommentReply struct {
	ReplyID			int
	Content			string
	Time			string
	Author			string
	Self			bool
}
