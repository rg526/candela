package cdmodel

type Comment struct {
	CommentID		int
	Content			string
	Time			string
	Author			string
	Self			bool
	Replies			[]CommentReply
}

type CommentReply struct {
	ReplyID			int
	Content			string
	Time			string
	Author			string
	Self			bool
}
