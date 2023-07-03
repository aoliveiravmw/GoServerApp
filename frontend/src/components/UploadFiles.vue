<template>
 <div class="hello">
    <input type="file" @change="OnFileSelected">
    <button @click="onUpload">Upload</button>
 </div>
</template>

<script>
import axios from 'axios'

export default {
    name: 'UploadFiles',
    data () {
        return {
           selectedFile: null
        }
    },
    methods: {
        OnFileSelected(event) {
            this.selectedFile = event.target.files[0]
        },
        onUpload() {
            const fd = new FormData();
            fd.append('picture', this.selectedFile)
            axios.post('http://192.168.223.23:8081/file', fd, {
                onUploadProgress: uploadEvent => {
                    console.log('Upload Progress: ' + Math.round(uploadEvent.loaded / uploadEvent.total *100) + '%')
                }
            })
            .then(res => {
                console.log(res)
            })
        }
    }
}
</script>


