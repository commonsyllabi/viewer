let init = () => {
    Vue.createApp({
        data() {
            return {
                cartridge: null,
                manifest: null,
                resources: null,
                log: "",
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
                    if (res.ok) {
                        return res.json()
                    } else {
                        console.error(res.err)
                        this.log = `internal server error on upload: ${res.err}`
                    }
                }).then(data => {
                    console.log(data);
                    this.log = `loaded ${this.cartridge.name}`
                    this.manifest = JSON.parse(data.data)
                    this.resources = JSON.parse(data.resources)
                }).catch(err => {
                    console.error(err);
                })
            },
            getFile(_evt, _id) {
                fetch(`/api/file/${_id}?cartridge=${this.cartridge.name}`, {
                    method: 'GET'
                })
                    .then(res => { return res.json() })
                    .then(body => {
                        //-- assign to iframe
                        _evt.target.nextElementSibling.children[0].src = body.path 
                    })
                    .catch(err => console.error(err));
            }
        },
        mounted() {

        }
    }).mount('#app')
}