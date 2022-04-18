let init = () => {
    Vue.createApp({
        data() {
            return {
                cartridge : null
            }
        },
        methods: {
            upload() {
                let form = document.getElementById("upload-form")
                let formData = new FormData(form)
            
                if(formData.get("cartridge").name == ""){
                    console.warn("can't submit an empty cartridge!");
                    return
                }
            
                fetch("/api/upload",{
                    method: 'POST',
                    body: formData
                }).then(res => {
                    return res.json()
                }).then(data => {
                    console.log(data);
                }).catch(err => {
                    console.error(err);
                })
            }
        },
        mounted() {
            this.cartridge = DATA
        }
    }).mount('#app')
}