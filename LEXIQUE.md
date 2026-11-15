# Lexique

Les termes et abréviations employés dans les slides, module par module, en une phrase chacun.

## 1. Organisation du code

- **Package** : Un dossier de fichiers .go qui partagent le même nom de package. L'unité de visibilité : majuscule = exporté, minuscule = privé au package.
- **Module** : Un ensemble de packages versionné ensemble, défini par un go.mod. L'unité de publication et de dépendance.
- **internal/** : Dossier spécial : ses packages ne sont importables que depuis le code du même projet. Le compilateur refuse le reste.
- **Type métier (domaine)** : Une struct qui représente une notion du métier (Order) avec ses règles. Champs privés, construction validée, transitions contrôlées.
- **DTO (Data Transfer Object)** : Une struct qui ne sert qu'à transporter des données vers ou depuis l'extérieur (JSON, base). Tout public, zéro logique.
- **Interface** : Un contrat : la liste des méthodes qu'un type doit avoir. En Go un type la satisfait sans le déclarer (implémentation implicite).
- **Valeur zéro (ou vierge)** : Ce que vaut une variable déclarée sans initialisation : 0, "", nil, ou une struct dont tous les champs sont à zéro. Un bon type est utilisable tel quel.
- **Groupes de sérialisation (sheriff)** : Un tag groups:"list,detail" sur les champs, et une bibliothèque tierce (github.com/liip/sheriff) qui ne garde que les champs du groupe demandé avant json.Marshal. Une vue par route (liste, détail) sans struct supplémentaire.
- **reflect** : Le package qui permet d'examiner et de manipuler les types à l'exécution (champs, tags, méthodes par leur nom). Les bibliothèques génériques en vivent, le code métier doit s'en passer.
- **unsafe** : Le package qui contourne le système de types : conversions de pointeurs arbitraires, accès mémoire direct. Sizeof et Offsetof mesurent sans danger, le reste est à proscrire.
- **Functional options** : Un constructeur New(opts ...Option) où chaque option est une fonction qui modifie la struct. Permet d'ajouter un réglage sans casser les appels existants.

## 2. Gestion des erreurs

- **Sentinelle** : Une erreur déclarée une fois comme variable (var ErrNotFound = errors.New(...)) et comparée avec errors.Is. Pour dire « laquelle », pas « avec quels détails ».
- **Erreur typée** : Une struct qui implémente error et transporte des champs (Field, Step...). On la récupère avec errors.As.
- **Enrober (wrap)** : Créer une nouvelle erreur qui contient l'ancienne : fmt.Errorf("contexte: %w", err). La chaîne reste inspectable.
- **%w** : Le verbe de fmt.Errorf qui enrobe. Avec %v l'erreur d'origine devient du texte et errors.Is ne la retrouve plus.
- **Chaîne d'erreurs** : La suite des erreurs enrobées les unes dans les autres, de la plus externe à la cause. errors.Is et errors.As la parcourent.
- **Unwrap** : La méthode qu'une erreur typée implémente pour exposer celle qu'elle contient. C'est elle qui rend la chaîne traversable.
- **errors.Join** : Regroupe plusieurs erreurs en une seule (une étape qui échoue puis une compensation qui échoue aussi).
- **Frontière** : L'endroit où une erreur Go devient autre chose : un code HTTP, un message utilisateur. Une seule fonction s'en charge (writeError).

## 3. Logging structuré

- **slog** : Le package de logging structuré de la bibliothèque standard (log/slog, depuis Go 1.21).
- **logrus** : Bibliothèque tierce de log structuré, très répandue avant slog. Compatible avec l'API du package log, en maintenance depuis 2020.
- **Log structuré** : Un log fait de paires clé = valeur plutôt que d'une phrase. Une machine peut le filtrer, un humain peut le lire.
- **Handler** : Le composant de slog qui écrit réellement : texte, JSON, ou le vôtre. Le Logger lui transmet des Records.
- **Record** : Un événement de log : horodatage, niveau, message, attributs. Ce que le Handler reçoit.
- **Attribut** : Une paire clé / valeur attachée à un log ("order_id", id). With() en attache à un logger pour tous ses logs suivants.
- **Niveau** : Debug, Info, Warn, Error. Le handler ignore ce qui est sous son niveau minimum.
- **request_id** : Un identifiant unique par requête HTTP, mis dans chaque log de la requête pour pouvoir les retrouver ensemble.
- **LogValuer** : Une interface : un type qui sait se décrire pour les logs (LogValue()). Order choisit ce qu'il montre, sans exposer sa struct.

## 4. Injection de dépendances

- **Injection de dépendances (DI)** : Donner à un composant ce dont il a besoin de l'extérieur, au lieu qu'il le construise lui-même. fx le fait automatiquement.
- **Graphe de dépendances** : L'ensemble des constructeurs reliés par leurs types : A a besoin de B qui a besoin de C. fx le résout par types.
- **Constructeur** : Une fonction qui renvoie un type (et éventuellement une erreur). Ce qu'on donne à fx.Provide.
- **Provide** : Déclare à fx comment fabriquer un type. Rien n'est exécuté tant que personne n'en a besoin.
- **Invoke** : Demande à fx d'appeler une fonction au démarrage, en lui fournissant ses paramètres. C'est ce qui tire le graphe.
- **Lifecycle / hook** : Le cycle de vie de l'application. Un hook OnStart / OnStop s'exécute au démarrage et à l'arrêt, dans l'ordre inverse.
- **fx.Module** : Un regroupement nommé de Provide et d'Invoke. Un module par zone fonctionnelle.
- **fx.Annotate, fx.As** : Fournir un type concret sous une interface, ou le nommer quand il y a deux instances du même type.
- **fx.Supply, fx.Options** : Supply fournit une valeur déjà construite au graphe. Options regroupe des Provide et Invoke sans nom (Module est la version nommée).
- **fx.Decorate, fx.Replace, fx.Populate** : Decorate modifie une valeur du graphe pour la suite, Replace la substitue (tests), Populate la copie dans une variable à soi.
- **fx.In / fx.Out** : Des structs qui remplacent une longue liste de paramètres ou de résultats de constructeur.

## 5. Modules privés

- **Semver** : Versionnage sémantique : vMAJEUR.MINEUR.CORRECTIF. Go s'en sert pour résoudre les versions et détecter les ruptures.
- **Tag** : Une étiquette git posée sur un commit (v0.1.0). Pour Go, un tag = une version publiée du module.
- **go.mod / go.sum** : go.mod déclare le module et ses dépendances, go.sum enregistre leurs empreintes pour vérifier qu'elles n'ont pas changé.
- **go.work** : Un fichier d'espace de travail qui fait pointer plusieurs modules locaux les uns vers les autres, sans toucher aux go.mod. Pour le poste de dev, pas pour le dépôt.
- **replace** : Une directive de go.mod qui remplace une dépendance par un chemin local ou un fork. Facile à commiter par erreur.
- **Proxy de modules** : Le serveur (proxy.golang.org par défaut) par lequel go télécharge les modules et vérifie leurs empreintes.
- **GOPRIVATE** : Variable d'environnement : les modules qui correspondent ne passent ni par le proxy public ni par la base d'empreintes. Obligatoire pour un dépôt privé.
- **GONOSUMDB / GONOPROXY** : Les deux variables que GOPRIVATE règle d'un coup. On peut les régler séparément si besoin.

## 6. context.Context

- **context.Context** : Un objet qui voyage avec un appel et porte trois choses : une annulation, un délai, des valeurs. Toujours en premier paramètre.
- **Annulation** : Le signal que le travail n'est plus attendu. ctx.Done() se ferme, ctx.Err() dit pourquoi.
- **Deadline / timeout** : Un instant (deadline) ou une durée (timeout) après quoi le contexte s'annule tout seul.
- **Parent / enfant** : Un contexte dérivé d'un autre. Annuler le parent annule tous les enfants, jamais l'inverse.
- **Cause** : L'erreur précise de l'annulation (WithCancelCause, WithTimeoutCause), lue avec context.Cause(ctx). Plus parlante que ctx.Err().
- **WithoutCancel** : Un contexte qui garde les valeurs du parent mais ne s'annule plus avec lui. Pour une compensation qui doit aller au bout.
- **Propagation** : Passer le contexte reçu à tout ce qu'on appelle, pour que l'annulation descende jusqu'au bout.
- **Valeur de contexte** : Une donnée attachée au contexte (request_id, logger). Pour ce qui traverse les couches, jamais pour les paramètres métier.

## 7. Fuites de goroutines

- **Goroutine** : Un fil d'exécution léger géré par Go. On en lance des milliers, mais chacune doit avoir une fin.
- **Fuite de goroutines** : Une goroutine qui ne se terminera jamais parce qu'elle attend quelque chose qui ne viendra pas. La mémoire monte, le service s'étouffe.
- **pprof** : L'outil de profilage de Go. Sur /debug/pprof/goroutine il liste toutes les goroutines vivantes et où elles sont bloquées.
- **Profil goroutineleak** : Nouveau profil pprof de Go 1.27 (/debug/pprof/goroutineleak) : ne liste que les goroutines bloquées sur un canal ou un verrou devenu inatteignable, donc les fuites. Aveugle aux attentes réseau et aux Sleep.
- **Pile (stack trace)** : La liste des fonctions en cours d'appel pour une goroutine. C'est ce qu'on lit pour savoir où elle attend.
- **goleak** : Une bibliothèque de test (go.uber.org/goleak) qui fait échouer un test s'il reste des goroutines à la fin.
- **Race detector** : go test -race : détecte deux goroutines qui touchent la même mémoire sans synchronisation.
- **select** : L'instruction qui attend sur plusieurs canaux à la fois. Avec une branche sur ctx.Done(), une attente devient annulable.
- **Channel (canal)** : Une file de communication entre goroutines. Bufferisé ou non, il bloque quand il est plein ou vide.

## 8. Borner la concurrence

- **Pool de workers** : Un nombre fixe de goroutines qui lisent des tâches dans une file. On borne la concurrence au lieu de lancer une goroutine par requête.
- **File bornée** : Une file d'attente de taille fixe. Pleine, elle refuse : c'est le délestage.
- **Délestage (back-pressure)** : Refuser explicitement du travail quand on sature (503 + Retry-After) plutôt que de l'accepter et s'effondrer.
- **Sémaphore** : Un compteur de jetons : au plus N appels en même temps vers une ressource. x/sync/semaphore.Weighted.
- **Mutex / RWMutex** : Un verrou. RWMutex laisse plusieurs lecteurs en même temps mais un seul écrivain.
- **Data race** : Deux goroutines accèdent à la même donnée, dont une en écriture, sans verrou. Résultat indéfini, détecté par -race.
- **sync.Once** : Garantit qu'une fonction ne s'exécute qu'une fois, même appelée depuis plusieurs goroutines.
- **Retry-After** : En-tête HTTP qui dit au client dans combien de temps réessayer. Accompagne un 503.

## 9. Arrêt propre

- **Arrêt propre (graceful shutdown)** : Sur un signal d'arrêt : ne plus rien accepter, finir ce qui est en cours, puis sortir. Sans perdre de travail en silence.
- **SIGTERM / SIGINT** : Les signaux envoyés par l'orchestrateur (Kubernetes, systemd) ou par Ctrl-C pour demander l'arrêt. fx.Run les écoute.
- **Drainer** : Laisser les workers finir les tâches déjà prises, sans en accepter de nouvelles.
- **StopTimeout** : Le temps maximal accordé à l'arrêt. Au-delà, on annule ce qui reste. fx.StopTimeout côté fx, le même délai moins une marge côté pool.
- **Arrêt forcé (hard stop)** : L'annulation des commandes encore en cours quand le délai est écoulé, avec une cause explicite pour qu'on le voie dans les logs.
- **context.AfterFunc** : Enregistre une fonction à appeler quand un contexte s'annule. Pratique pour tracer chaque annulation.
- **kill -9 (SIGKILL)** : L'arrêt qu'on ne peut pas intercepter. C'est ce qui arrive quand le délai de l'orchestrateur est dépassé : tout ce qui n'est pas fini est perdu.

## 10. Tests

- **Test unitaire** : Un test qui vérifie un composant seul, ses dépendances remplacées par des doublures.
- **Table-driven test** : Un test écrit comme une liste de cas (entrée, attendu) parcourue en boucle. Le style idiomatique en Go.
- **assert vs require** : testify : assert note l'échec et continue, require arrête le test tout de suite. require pour les préconditions.
- **Suite** : testify/suite : une struct qui regroupe des tests avec un SetupTest exécuté avant chacun. Évite de répéter le câblage.
- **Fake / mock / stub** : Fake : une vraie implémentation simplifiée (en mémoire). Mock : vérifie les appels reçus. Stub : renvoie des valeurs fixes.
- **mockery** : L'outil qui génère les mocks testify à partir des interfaces. Configuré dans .mockery.yml, v3.
- **httptest** : Le package standard pour tester un handler HTTP sans réseau : httptest.NewRequest, httptest.NewRecorder.
- **fxtest, fx.Replace** : fxtest.New démarre une app fx dans un test et échoue proprement. fx.Replace substitue un composant du graphe.
- **synctest** : testing/synctest : exécute un test dans un temps virtuel. Un timeout de 5 s passe en quelques millisecondes.
- **Couverture** : La part du code exécutée par les tests (go test -cover). Un indicateur, pas un objectif.

## 11. Débugger en autonomie

- **Panic** : L'arrêt brutal d'une goroutine sur une erreur irrécupérable (nil pointer, index hors limites). Non rattrapé, il tue le processus.
- **Trace de pile** : Ce que Go imprime à un panic : pour chaque goroutine, la liste des appels en cours. Se lit de haut en bas, première ligne de votre code = la bonne.
- **Deadlock** : Des goroutines qui s'attendent mutuellement pour toujours. Go le détecte quand tout le monde dort : « all goroutines are asleep ».
- **SIGQUIT** : Le signal (Ctrl-\) qui fait imprimer à un programme Go toutes ses piles sans l'arrêter proprement. Pour un service qui semble figé.
- **Delve (dlv)** : Le débogueur de Go. dlv debug, dlv test, dlv attach, dlv exec.
- **Point d'arrêt (breakpoint)** : Une ligne où le débogueur s'arrête. Conditionnel : seulement si une expression est vraie (attempt == 3).
- **attach** : Brancher le débogueur sur un processus déjà en cours d'exécution, par son PID.
- **launch.json** : La configuration de débogage de VS Code : quoi lancer, avec quels arguments et variables d'environnement.

## 12. Performance

- **Profil** : Une mesure de où le programme passe son temps (CPU) ou alloue sa mémoire (mémoire). Produit par pprof, à lire avant d'optimiser.
- **b.Loop** : La boucle de benchmark moderne (Go 1.24) : for b.Loop() { ... }. Remplace for i := 0; i < b.N; i++ et empêche le compilateur d'optimiser le code mesuré.
- **Benchmark** : Une fonction BenchmarkXxx exécutée en boucle par go test -bench. Donne ns/op, B/op et allocs/op.
- **Allocation** : Une réservation de mémoire sur le tas, que le ramasse-miettes devra libérer. Moins on en fait, moins on paie.
- **benchstat** : L'outil qui compare deux séries de benchmarks et dit si la différence est significative.
- **Alignement** : La façon dont le compilateur place les champs d'une struct en mémoire. Leur ordre change la taille de la struct.

## 13. Sécurité

- **Durcissement (hardening)** : Fermer les portes qu'un service laisse ouvertes par défaut : timeouts, taille des corps, en-têtes, origines autorisées, récupération des panics.
- **Timeouts du http.Server** : ReadHeaderTimeout, ReadTimeout, WriteTimeout, IdleTimeout. Sans eux, un client lent ou malveillant garde une connexion pour toujours.
- **Path traversal** : Une attaque où un chemin fourni par le client (../../etc/passwd) sort du dossier prévu.
- **os.OpenRoot** : API Go 1.24 : ouvre des fichiers relativement à un dossier racine, sans possibilité d'en sortir.
- **CORS** : Les en-têtes qui disent à un navigateur quelles origines (sites) peuvent appeler l'API. En liste blanche, jamais *.
- **MaxBytesReader** : Limite la taille d'un corps de requête. Au-delà, lecture refusée et 413.
- **nosniff** : L'en-tête X-Content-Type-Options: nosniff : le navigateur ne devine pas le type d'un fichier, il croit le Content-Type.
- **Recover** : Un middleware qui rattrape un panic dans un handler, logue, et répond 500 au lieu de laisser tomber la connexion.
- **govulncheck** : L'outil officiel qui liste les vulnérabilités connues des dépendances réellement appelées par votre code.

# Pour aller plus loin

Les ressources citées en bas des slides.

- [Golang Functional Options Pattern (golang.cafe)](https://golang.cafe/blog/golang-functional-options-pattern)
- [Functional options for friendly APIs, Dave Cheney](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
- [sheriff sur pkg.go.dev](https://pkg.go.dev/github.com/liip/sheriff/v2)
- [Go Proverbs, Rob Pike : « The bigger the interface, the weaker the abstraction »](https://go-proverbs.github.io/)
- [Working with Errors in Go 1.13 (blog Go)](https://go.dev/blog/go1.13-errors)
- [Errors (Go by Example)](https://gobyexample.com/errors)
- [Joining errors in Golang (GeeksforGeeks)](https://www.geeksforgeeks.org/go-language/joining-errors-in-golang/)
- [Package errors (pkg.go.dev)](https://pkg.go.dev/errors)
- [Error handling and Go (blog Go)](https://go.dev/blog/error-handling-and-go)
- [logrus, le README et son statut](https://github.com/sirupsen/logrus)
- [slog-logrus, le pont entre les deux](https://github.com/samber/slog-logrus)
- [Structured Logging with slog (blog Go)](https://go.dev/blog/slog)
- [Logging in Go with slog, le guide Better Stack](https://betterstack.com/community/guides/logging/logging-in-go/)
- [Contextual Logging in Go with slog (Better Stack)](https://betterstack.com/community/guides/logging/golang-contextual-logging/)
- [Publishing a module (go.dev)](https://go.dev/doc/modules/publishing)
- [Tutorial : Getting started with multi-module workspaces](https://go.dev/doc/tutorial/workspaces)
- [Go Modules Reference, Private modules](https://go.dev/ref/mod#private-modules)
- [fx : Get started](https://uber-go.github.io/fx/get-started/index.html)
- [fx : Modules](https://uber-go.github.io/fx/modules.html)
- [fx : Annotations](https://uber-go.github.io/fx/annotate.html)
- [fx : Parameter Objects et Result Objects](https://uber-go.github.io/fx/parameter-objects.html)
- [fx : Application lifecycle](https://uber-go.github.io/fx/lifecycle.html)
- [fx : Supply, Decorate, Replace, Populate (pkg.go.dev)](https://pkg.go.dev/go.uber.org/fx)
- [Go Concurrency Patterns : Context (blog Go)](https://go.dev/blog/context)
- [Package context (pkg.go.dev)](https://pkg.go.dev/context)
- [Contexts and structs (blog Go)](https://go.dev/blog/context-and-structs)
- [Goroutine Leaks, The Forgotten Sender (Ardan Labs)](https://www.ardanlabs.com/blog/2018/11/goroutine-leaks-the-forgotten-sender.html)
- [Go 1.27 : le profil goroutineleak (Ardan Labs)](https://www.ardanlabs.com/news/2026/go-1-27-goroutine-leak-profiling-why-goroutine-leaks-are-difficult-to-catch)
- [Never start a goroutine without knowing how it will stop, Dave Cheney](https://dave.cheney.net/2016/12/22/never-start-a-goroutine-without-knowing-how-it-will-stop)
- [uber-go/goleak](https://github.com/uber-go/goleak)
- [Worker Pools (Go by Example)](https://gobyexample.com/worker-pools)
- [Go Concurrency Patterns : Pipelines and cancellation (blog Go)](https://go.dev/blog/pipelines)
- [golang.org/x/sync/semaphore](https://pkg.go.dev/golang.org/x/sync/semaphore)
- [Data Race Detector (go.dev)](https://go.dev/doc/articles/race_detector)
- [Graceful Shutdown in Go : Practical Patterns (VictoriaMetrics)](https://victoriametrics.com/blog/go-graceful-shutdown/)
- [http.Server.Shutdown (pkg.go.dev)](https://pkg.go.dev/net/http#Server.Shutdown)
- [Table Driven Tests (Go wiki)](https://go.dev/wiki/TableDrivenTests)
- [Documentation mockery v3](https://vektra.github.io/mockery/latest/)
- [testify/suite (pkg.go.dev)](https://pkg.go.dev/github.com/stretchr/testify/suite)
- [net/http/httptest (pkg.go.dev)](https://pkg.go.dev/net/http/httptest)
- [fx : Testing](https://pkg.go.dev/go.uber.org/fx/fxtest)
- [Testing concurrent code with testing/synctest (blog Go)](https://go.dev/blog/testing-time)
- [Diagnostics (go.dev)](https://go.dev/doc/diagnostics)
- [Delve : commandes du débogueur](https://github.com/go-delve/delve/tree/master/Documentation/cli)
- [vscode-go : Debugging (wiki)](https://github.com/golang/vscode-go/wiki/debugging)
- [GoLand : Debug Go code](https://www.jetbrains.com/help/go/debugging-code.html)
- [Profiling Go Programs (blog Go)](https://go.dev/blog/pprof)
- [High Performance Go Workshop, Dave Cheney](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)
- [More predictable benchmarking with testing.B.Loop (blog Go)](https://go.dev/blog/testing-b-loop)
- [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
- [Traversal-resistant file APIs (blog Go)](https://go.dev/blog/osroot)
- [The complete guide to Go net/http timeouts (Cloudflare)](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
- [Tutorial : Find and fix vulnerable dependencies with govulncheck](https://go.dev/doc/tutorial/govulncheck)
- [OWASP Go Secure Coding Practices](https://github.com/OWASP/Go-SCP)
