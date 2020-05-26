# Supply chain Management

Nous allons réaliser la simulation d’une application de Supply Chain management,
Nous allons pour cela simuler les différentes étapes se trouvant avant la réalisation d’un produit.

Chaque produit passe  par une phase de fabrication, une phase de conditionnement, une phase de transport, une phase de stockage, puis une dernière phase de transport pour aboutir à la phase de commercialisation par les détaillants.

Choix de l’implémentation

Le choix de l’implémentation et du langage étant laissé à notre bon-vouloir, nous avons opté pour le langage Go et l’implémentation des espaces de Tuples [GoSpace](https://github.com/pSpaces/goSpace).

Ces choix ont été motivés par la facilité de programmation concurrente de go et la présence de l’ensemble des primitives d’accès de Linda dans l’implémentation Gospace.

