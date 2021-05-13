package objects

class Node {
    def uid
    def name
    def labels
    Map<String, String> annotations
    def internalIps
    def externalIps
    def containerRuntimeVersion
    def kernelVersion
    def osImage
    def kubeletVersion
    def kubeProxyVersion
}
