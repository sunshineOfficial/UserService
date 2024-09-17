package kafka

type ProducerOptions struct {
	BatchSize  int
	BatchBytes int64
	Async      bool
}

type ProducerOption func(p ProducerOptions) ProducerOptions

func WithBatchSize(size int) ProducerOption {
	return func(p ProducerOptions) ProducerOptions {
		p.BatchSize = size
		return p
	}
}

func WithBatchBytes(bytes int64) ProducerOption {
	return func(p ProducerOptions) ProducerOptions {
		p.BatchBytes = bytes
		return p
	}
}

func ProduceAsync() ProducerOption {
	return func(p ProducerOptions) ProducerOptions {
		p.Async = true
		return p
	}
}
