# Controller

What about:
* We create a controller that will read CRD for defining the gateway on a namespace
  * the controller will create
    * daemon on all nodes
    * headless service for reaching the daemon set from the node
    * transform CRD into Config map that will be later on injected into the pod
    * If it has dlq, then an option to allow the CRD create one

# How to ensure resilience
* WAL
  * Implement something like checkpoints on files
    * then I could analyzed the file if the pod goes down
    * the files are replicated across nodes, then we could easily catch up with it
* Implementing a control system on the influx?