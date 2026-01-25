package router

func Link(a, b *Router, cost float64) {
	a.UpdateNeighbor(b.ID(), cost)
	b.UpdateNeighbor(a.ID(), cost)
}
