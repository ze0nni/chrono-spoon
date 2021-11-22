VueWS = (function() {

const VueWS = {
        install(Vue, connectionUrl) {
                const wsConnection = WsConnection(connectionUrl)
                Vue.prototype.$wsocket = wsConnection.wsocket

                Vue.mixin({
                        created,
                        beforeDestroy
                })

                function created() {
                        wsConnection.listen(this)
                }
                
                function beforeDestroy() {
                        wsConnection.unlisten(this)
                }
        }
}

return VueWS

function WsConnection(url) {
        let isConnected = false
        const listenedComponents = {}

        const ws = new WebSocket(url)
        ws.onopen = function() {
                isConnected = true

                for (component of Object.values(listenedComponents)) {
                        const connected = component.$options.webSockets.connected;
                        if (null != connected) {
                                connected.call(component)
                        }
                }
        }
        ws.onerror = function() {
                isConnected = false
        }
        ws.onclose = function() {
                isConnected = false
                for (component of Object.values(listenedComponents)) {
                        const disconnected = component.$options.webSockets.disconnected;
                        if (null != disconnected) {
                                disconnected.call(component)
                        }
                }
        }
        ws.onmessage = function(event) {
                const msg = JSON.parse(event.data)
                const command = msg.command
                for (component of Object.values(listenedComponents)) {
                        const wsOptions = component.$options.webSockets
                        if (wsOptions.command) {
                                const reciever = wsOptions.command[command]
                                if (reciever) {
                                        reciever.call(component, msg)
                                }
                        }
                }
        }

        function sendMessage(message) {
                if (!isConnected) {
                        throw new Error("Ws not ready")
                }
                ws.send(JSON.stringify(message))
        }

        const wsocket = {
                send: sendMessage
        }

        function listen(component) {
                const wsOptions = component.$options.webSockets
                if (null == wsOptions) {
                        return
                }
                listenedComponents[component._uid] = component
                
                if (isConnected && wsOptions.connected) {
                        wsOptions.connected.call(component)
                }
        }

        function unlisten(component) {
                delete listenedComponents[component._uid]
        }

        return {
                wsocket,
                listen,
                unlisten
        }
}

})()