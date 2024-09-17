package kafka

type ConsumerOptions struct {
	GroupId       string
	TopicName     string
	Partition     int
	QueueCapacity int
	MinBytes      int
	MaxBytes      int
	StartOffset   int64
}

type ConsumerOption func(p ConsumerOptions) ConsumerOptions

func WithConsumerGroup(name string) ConsumerOption {
	return func(p ConsumerOptions) ConsumerOptions {
		p.GroupId = name
		return p
	}
}

func WithTopic(name string) ConsumerOption {
	return func(p ConsumerOptions) ConsumerOptions {
		p.TopicName = name
		return p
	}
}

func WithPartition(id int) ConsumerOption {
	return func(p ConsumerOptions) ConsumerOptions {
		p.Partition = id
		return p
	}
}

func WithQueue(capacity int) ConsumerOption {
	return func(p ConsumerOptions) ConsumerOptions {
		p.QueueCapacity = capacity
		return p
	}
}

func WithMinBytes(min int) ConsumerOption {
	return func(p ConsumerOptions) ConsumerOptions {
		p.MinBytes = min
		return p
	}
}

func WithMaxBytes(max int) ConsumerOption {
	return func(p ConsumerOptions) ConsumerOptions {
		p.MaxBytes = max
		return p
	}
}

func WithOffset(start int64) ConsumerOption {
	return func(p ConsumerOptions) ConsumerOptions {
		p.StartOffset = start
		return p
	}
}
