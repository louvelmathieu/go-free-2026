# Étape 9 : arrêt propre

- `Pool.Stop(ctx)` : plus d'admission (`ErrStopping` → 503 « redémarrage »),
  fermeture de la file sous lock, les workers drainent ce qui reste,
  attente bornée par le délai reçu moins une marge de 500 ms.
- Au-delà du délai : `cancelHardStop(ErrShuttingDown)`. Chaque job relie
  ce signal à son propre contexte avec `context.AfterFunc` (un contexte
  n'a qu'un parent). Les commandes interrompues finissent `cancelled`
  avec une cause explicite, aucune n'est perdue en silence.
- `fx.StopTimeout(cfg.StopTimeout)` (`APP_STOP_TIMEOUT`, 20 s par défaut).
  La config est lue avant `fx.New` et réinjectée par `fx.Supply`.
- Ordre d'arrêt (inverse des OnStart) : admin, HTTP (`Shutdown` attend
  les requêtes en cours), puis workers.
- Démo : 20 POST puis Ctrl-C. Délai 20 s : 20 terminées. Délai 2 s :
  4 terminées, 16 annulées « arrêt du service avant la fin de la commande ».
