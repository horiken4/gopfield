package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Axon chan float32

type Neuron struct {
	id        int
	weights   map[*Neuron]float32 // Weights for input axons
	inAxons   map[*Neuron]Axon    // Input axons
	outAxons  map[*Neuron]Axon    // Output axons
	v         float32             // Electrical potential
	vChan     chan float32        // Channel for stimulate electrical potential
	th        float32             // Threshold
	trainMode bool                // If true neuron conducts training
}

func NewNeuron() *Neuron {
	return &Neuron{
		weights:  map[*Neuron]float32{},
		inAxons:  map[*Neuron]Axon{},
		outAxons: map[*Neuron]Axon{},
		vChan:    make(chan float32, 1),
	}
}

// Connect connects axon between n and neuron
func (n *Neuron) Connect(neuron *Neuron) error {
	var (
		axon Axon
		ok   bool
	)

	fmt.Printf("Connectting between neuron %v to neuron %v\n", n.id, neuron.id)

	// Inbound axon (from neuron to n)
	axon, ok = n.inAxons[neuron]
	if !ok {
		n.inAxons[neuron] = make(chan float32, 1)
		axon = n.inAxons[neuron]
	}
	if _, ok := neuron.outAxons[n]; ok {
		return errors.New(fmt.Sprintf("Connection from neuron %v to neuron %v has already existed", neuron.id, n.id))
	}
	neuron.outAxons[n] = axon

	// Outbound axon (from n to neuron)
	axon, ok = neuron.inAxons[n]
	if !ok {
		neuron.inAxons[n] = make(chan float32, 1)
		axon = neuron.inAxons[n]
	}
	if _, ok := n.outAxons[neuron]; ok {
		return errors.New(fmt.Sprintf("Connection from neuron %v to neuron %v has already existed", n.id, neuron.id))
	}
	n.outAxons[neuron] = axon

	return nil
}

// Feed feeds -1 or 1 to this neuron
func (n *Neuron) Feed(v float32) error {
	if v != -1 && v != 1 {
		return errors.New("Invalid feed value")
	}

	n.vChan <- v

	// TODO: Add channel for notification feed has finisied

	return nil
}

// Run runs neuron processes
func (n *Neuron) Run(iter int, finish chan bool) {
	go n.run(iter, finish)
}

func (n *Neuron) run(iter int, finish chan bool) {
	fmt.Printf("neuron %v : start\n", n.id)

	if n.trainMode {
		// Initialize weights to 0
		for neuron := range n.weights {
			n.weights[neuron] = 0
		}

		for it := 0; it < iter; it++ {
			// Initialize electrical potential
			n.v = <-n.vChan
			for _, axon := range n.outAxons {
				axon <- n.v
			}

			// Update weights by Hebb's rule
			for neuron, axon := range n.inAxons {
				n.weights[neuron] += n.v * <-axon
			}
		}
	} else {

		// Initialize electrical potential
		n.v = <-n.vChan

		for it := 0; it < iter; it++ {
			fmt.Println(n.id, "neuron %v : it =", n.id, it)

			for neuron, axon := range n.outAxons {
				fmt.Printf("neuron %v : output to %v before\n", n.id, neuron.id)
				axon <- n.v
				fmt.Printf("neuron %v : output to %v after\n", n.id, neuron.id)
			}

			var s float32
			for neuron, axon := range n.inAxons {
				fmt.Printf("neuron %v : input from %v before\n", n.id, neuron.id)
				s += n.weights[neuron] * <-axon
				fmt.Printf("neuron %v : input from %v after\n", n.id, neuron.id)
			}

			s -= n.th

			var v float32 = -1
			if s > 0 {
				v = 1
			}
			n.v = v
		}
	}

	fmt.Printf("neuron %v : finish\n", n.id)

	finish <- true
}

type Hopfield struct {
	Neurons []*Neuron
}

func NewHopfield(numNeurons int) *Hopfield {
	// Make neurons
	neurons := make([]*Neuron, numNeurons)
	for i := 0; i < numNeurons; i++ {
		neurons[i] = NewNeuron()
		neurons[i].id = i
	}

	// Make axons
	for i := 0; i < numNeurons; i++ {
		for j := i + 1; j < numNeurons; j++ {
			err := neurons[i].Connect(neurons[j])
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
		}
	}

	return &Hopfield{
		Neurons: neurons,
	}
}

func (h *Hopfield) Energy() float32 {
	var (
		e1 float32
		e2 float32
	)

	for i := 0; i < len(h.Neurons); i++ {
		for j := 0; j < len(h.Neurons); j++ {
			u := h.Neurons[i]
			v := h.Neurons[j]
			e1 += u.weights[v] * u.v * v.v
		}
	}
	e1 /= -2

	for i := 0; i < len(h.Neurons); i++ {
		u := h.Neurons[i]
		e2 += u.th * u.v
	}

	return e1 + e2
}

func (h *Hopfield) Run(iter int) {
	finish := make(chan bool, len(h.Neurons))

	for _, neuron := range h.Neurons {
		neuron.Run(iter, finish)
	}

	for i := 0; i < len(h.Neurons); i++ {
		<-finish
	}
}

func (h *Hopfield) Print(cols int) {
	for i, neuron := range h.Neurons {
		if i != 0 && i%cols == 0 {
			fmt.Print("\n")
		}
		if neuron.v == -1 {
			fmt.Print("○")
		} else {
			fmt.Print("●")
		}
	}
	fmt.Print("\n")
}

func (h *Hopfield) Feed(pat []float32) error {
	if len(pat) != len(h.Neurons) {
		return errors.New("Pattern size must be same to hopfield")
	}
	for _, v := range pat {
		if v != -1 && v != 1 {
			return errors.New("Pattern value must be -1 or 1")
		}
	}

	for i, neuron := range h.Neurons {
		err := neuron.Feed(pat[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hopfield) FeedRandomly() error {
	rand.Seed(time.Now().Unix())
	pat := make([]float32, len(h.Neurons))
	for i := 0; i < len(pat); i++ {
		pat[i] = float32(1 - 2*rand.Intn(2))
	}
	err := h.Feed(pat)
	if err != nil {
		return err
	}

	return nil
}

func (h *Hopfield) Train(pats [][]float32) error {
	// Set all neurons as training mode
	for _, neuron := range h.Neurons {
		neuron.trainMode = true
	}

	for _, pat := range pats {
		if len(pat) != len(h.Neurons) {
			return errors.New("Pattern size must be same to hopfield")
		}
		for _, v := range pat {
			if v != -1 && v != 1 {
				return errors.New("Pattern value must be -1 or 1")
			}
		}
	}

	// Training

	finish := make(chan bool, len(h.Neurons))

	for _, neuron := range h.Neurons {
		neuron.Run(len(pats), finish)
	}

	for _, pat := range pats {
		for i, v := range pat {
			h.Neurons[i].Feed(v)
		}
	}

	for i := 0; i < len(h.Neurons); i++ {
		<-finish
	}

	// Set all neurons as default mode
	for _, neuron := range h.Neurons {
		neuron.trainMode = false
	}

	return nil
}

func (h *Hopfield) SetWeights(i, j int, w float32) {
	h.Neurons[i].weights[h.Neurons[j]] = w
	h.Neurons[j].weights[h.Neurons[i]] = w
}

func (h *Hopfield) SetThreshold(i int, th float32) {
	h.Neurons[i].th = th
}
