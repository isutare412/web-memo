package enum

type PublishState string

const (
	PublishStatePrivate   PublishState = "private"
	PublishStateShared    PublishState = "shared"
	PublishStatePublished PublishState = "published"
)

func (PublishState) Values() []string {
	return []string{
		string(PublishStatePrivate),
		string(PublishStateShared),
		string(PublishStatePublished),
	}
}
