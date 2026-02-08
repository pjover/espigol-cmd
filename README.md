# Espígol

Espígol és una eina de línia de comandes per automatitzar la gestió d'ajuts de la cooperativa: importa previsió de despeses des de CSV, guarda i classifica despeses i factures, i genera informes i resums de seguiment per soci i secció.


# Importar socis

Importa els socis des d'un fitxer CSV.

Utilitza la target del Makefile `importar-socis`. Passa `CSV` per utilitzar un fitxer personalitzat; si no s'especifica, per defecte s'utilitza `private/CSV/socis.csv`:

```bash
# usa per defecte private/CSV/socis.csv
make importar-socis

# especifica un fitxer CSV concret
make importar-socis CSV=~/Downloads/socis.csv
```
