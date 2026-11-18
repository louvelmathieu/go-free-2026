module gitlab.example.com/free/provisio

go 1.27

require (
	gitlab.example.com/free/freelib v0.1.0
	go.uber.org/fx v1.24.0
	go.uber.org/goleak v1.3.0
)

require (
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
)

// Dans ce kit de formation, gitlab.example.com n'existe pas : Go ne peut pas
// lire le go.mod de freelib@v0.1.0 pour construire le graphe des modules
// (il le fait même en mode workspace dès qu'une dépendance externe entre en
// jeu, ici fx). En entreprise, avec le tag publié sur le GitLab, go.work
// suffit et cette ligne ne doit pas exister.
replace gitlab.example.com/free/freelib => ../freelib
