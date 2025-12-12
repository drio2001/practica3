package taller

import (
	"math/rand"
	"sync"
	"time"
)

// =======================
// Fase del taller
// =======================

type FaseTaller struct {
	slots     chan struct{}
	maxCola   int
	enCola    int
	maxEnCola int
	mu        sync.RWMutex
	logger    *Logger
	faseID    Fase
	cfg       Config
	rnd       *rand.Rand
}

func NuevaFase(faseID Fase, capacidad int, maxCola int, cfg Config, logger *Logger, rnd *rand.Rand) *FaseTaller {
	return &FaseTaller{
		slots:   make(chan struct{}, capacidad),
		maxCola: maxCola,
		logger:  logger,
		faseID:  faseID,
		cfg:     cfg,
		rnd:     rnd,
	}
}

func (f *FaseTaller) Entrar(c Coche) func() {
	for {
		f.mu.Lock()
		if f.enCola < f.maxCola {
			f.enCola++
			if f.enCola > f.maxEnCola {
				f.maxEnCola = f.enCola
			}
			f.mu.Unlock()
			break
		}
		f.mu.Unlock()
		time.Sleep(5 * time.Millisecond)
	}

	f.slots <- struct{}{}
	f.mu.Lock()
	f.enCola--
	f.mu.Unlock()

	f.logger.Log(c, f.faseID, "ENTRA")

	return func() {
		f.logger.Log(c, f.faseID, "SALE")
		<-f.slots
	}
}

func (f *FaseTaller) Trabajar(c Coche) {
	base := baseDurationForCat(f.cfg, c.Categoria)
	time.Sleep(variedDuration(base, f.cfg, f.rnd))
}

// =======================
// Simulación RWMutex
// =======================

func RunSimulationRWMutex(cfg Config, logs bool) Stats {
	seed := cfg.SemillaRandom
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rnd := rand.New(rand.NewSource(seed))
	logger := NewLogger(logs)

	f1 := NuevaFase(FaseLlegada, cfg.NumPlazas, cfg.MaxColaFase1, cfg, logger, rnd)
	f2 := NuevaFase(FaseMecanico, cfg.NumMecanicos, cfg.MaxColaFase2, cfg, logger, rnd)
	f3 := NuevaFase(FaseLimpieza, cfg.NumPlazas, cfg.MaxColaFase3, cfg, logger, rnd)
	f4 := NuevaFase(FaseEntrega, cfg.NumPlazas, cfg.MaxColaFase4, cfg, logger, rnd)

	var coches []Coche
	id := 1
	for i := 0; i < cfg.CochesA; i++ {
		coches = append(coches, Coche{id, CatA, incidenciaFromCat(CatA)})
		id++
	}
	for i := 0; i < cfg.CochesB; i++ {
		coches = append(coches, Coche{id, CatB, incidenciaFromCat(CatB)})
		id++
	}
	for i := 0; i < cfg.CochesC; i++ {
		coches = append(coches, Coche{id, CatC, incidenciaFromCat(CatC)})
		id++
	}

	rnd.Shuffle(len(coches), func(i, j int) {
		coches[i], coches[j] = coches[j], coches[i]
	})

	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(coches))

	for _, car := range coches {
		c := car
		go func() {
			defer wg.Done()
			r1 := f1.Entrar(c)
			f1.Trabajar(c)
			r1()

			r2 := f2.Entrar(c)
			f2.Trabajar(c)
			r2()

			r3 := f3.Entrar(c)
			f3.Trabajar(c)
			r3()

			r4 := f4.Entrar(c)
			f4.Trabajar(c)
			r4()
		}()
	}

	wg.Wait()

	return Stats{
		TotalCoches: len(coches),
		Duracion:    time.Since(start),
		Estrategia:  "RWMutex",
	}
}
