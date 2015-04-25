# A Hopfield networks implementation by Go

Neuron update process is done by goroutine associated to each neuron. Namely one neuron is regarded as one goroutine. A connection (axon) between neurons is described by channel.
