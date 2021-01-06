## Laboratorio N°3 Sistemas distribuidos Grupo GCR - G
##### Integrantes:

| Nombre  |   Rol|
| ------------ | ------------ |
| Cristian Bernal R.  | 201773026-9   |
|  Raúl Álvarez C. |  201773010-2 |



#### Como correr el código en cada maquina
Dependiendo de la máquina, se deberá ejecutar el comando make correspondiente:

| Comando make | Máquina     | Funcionalidad |
|--------------|-------------|---------------|
| runc         | cualquiera | Cliente     |
| runb        | 10.10.28.81   | Broker      |
| rund0      | 10.10.28.82 | DNS 0   |
| rund1       | 10.10.28.83 | DNS 1    |
| rund2      | 10.10.28.84 | DNS 2    |
| runa         | cualquiera  | Admin    |


#### Funcionamiento del Código
Se usó el lenguaje Go junto a gRPC para los envios de datos entre los nodos y el cliente

#### Observaciones
Es necesario que el broker y los dns esten funcionando para que admin funcione, puesto que este al iniciarse se "registra" con el broker para hacer funcionar la consistencia Read your Writes.    
Usamos append, update, delete de comandos de administrador.  




#### Funcionamiento del sistema
Por parte del administrador, cuando se inicia este lo primero que realiza es contactar a broker para que le asigne una id (claramente, esto se hace totalmente automático). Una vez hecho esto, se podrán usar los comandos que son append, update y delete con los formatos especificados en la pauta. Cuando se realiza uno de estos comandos, admin contacta a broker para pedirle una ip de DNS, si es primera vez que lo hace broker le asignara una dirección ip aleatoria y además la almacenará en memoria junto con la id del admin que solicito la ip; en caso contrario que no sea la primera vez que admin solicite una ip de DNS, broker le dará la que le asignó previamente, todo esto para mantener la consistencia de Read your Writes, ya que si le volviéramos a asignar una ip random nuevamente, esta consistencia no necesariamente se respetaría.  Cuando ocurre el proceso de merge de DNS, el dns0 informa a broker de que ocurrió este proceso, por lo que el broker lo que hará es borrar las ip's que designo a cada admin, ya que como todos los DNS's tendrán los mismos archivos, la consistencia se respetaría en este caso cuando un admin edite en cualquiera de los DNS. Esto significa, que una vez ocurrido un merge, si un admin quiere realizar un nuevo cambio, broker procede nuevamente a realizar un asignado aleatorio de ip's de dns, por lo que en resumen, cada 5 minutos broker asigna aleatoriamente ip's de DNS a cada admin.

Por el lado del Cliente, este al pedir una ip correspondiente a un sitio web, contacta a Broker, donde este intenta buscar en cada DNS (en orden aleatorio), si no se logra encontrar el sitio en ningun DNS, se imprime en pantalla de Cliente un mensaje de que no se ha logrado encontrar aquel. Si es que existe tal sitio con su ip, en el primer DNS que logre encontrar el sitio, el Broker retornará a Cliente la ip de esta con su respectivo reloj correspondiente. Para mantener el Monotonic Reads, primero miramos si en un arreglo en  Cliente donde se almacenan los sitio web con sus ip's  y relojes respectivos, existe tal sitio que fue solicitado, se procede a comparar los relojes de vectores, por lo que : 
- Si el reloj de vector recibido por Broker es menor al que poseemos almacenado en Cliente, se imprime en pantalla la ip del sitio que se tenia previamente almacenado en el arreglo
- Si el reloj de vector recibido por Broker es mayor al que poseemos almacenado en Cliente, se sobreescribe el sitio con los nuevos datos recibidos y se imprime en pantalla la direccion nueva
- Si no ocurre ninguno de los dos casos anteriormente nombrados, osea, el reloj de vector es igual, se procede a imprimir la ip que se tenia previamente almacenado en el arreglo

Finalmente, si el sitio web no estaba en aquel arreglo mencionado anteriormente, se procede a agregarlo con su reloj y su ip e imprimirlo en pantalla.
