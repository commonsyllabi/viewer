let init = () => {
    Vue.createApp({
        data() {
            return {
                cartridge: null,
                manifest: null,
                resources: null,
                log: "",
                preview: ""
            }
        },
        methods: {
            upload() {
                let form = document.getElementById("upload-form")
                let formData = new FormData(form)

                if (formData.get("cartridge").name == "") {
                    console.warn("can't submit an empty cartridge!");
                    this.log = "can't submit an empty cartridge!"
                    return
                }

                this.cartridge = formData.get("cartridge")

                fetch("/api/upload", {
                    method: 'POST',
                    body: formData
                }).then(res => {
                    return res.json()
                }).then(data => {
                    console.log(data);
                    this.log = `loaded ${this.cartridge.name}`
                    this.manifest = JSON.parse(data.data)
                    this.resources = JSON.parse(data.resources)
                }).catch(err => {
                    console.error(err);
                })
            },
            getResource(_id, _type) {
                if(_type == "webcontent"){
                    this.getFile(_id)
                }else{
                    fetch(`/api/resource/${_id}?cartridge=${this.cartridge.name}`, {
                        method: 'GET'
                    }).then(res => {
                        return res.json()
                    }).then(data => {
                        console.log(data);
                    }).catch(err => {
                        console.error(err)
                    })
                }
            },
            getFile(_id) {
                fetch(`/api/file/${_id}?cartridge=${this.cartridge.name}`, {
                    method: 'GET'
                })
                    .then(res => { return res.json() })
                    .then(body => { this.preview = body.path })
                    .catch(err => console.error(err));

            }
        },
        mounted() {
            
        }
    }).mount('#app')
}