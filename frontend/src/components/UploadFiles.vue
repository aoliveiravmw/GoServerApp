<template>
 <div class="hello">
    <table>
      <tr>
        <th>Choose a file to Upload</th>
        <td><input type="file" @change="OnFileSelected"></td>
        <td><button @click="onUpload">Upload</button></td>
      </tr>
      <tr><div>{{ uploadOutput[uploadOutput.length - 1] }}</div></tr>
    </table>
 </div>
</template>

<script>
import axios from 'axios'

export default {
    name: 'UploadFiles',
    data () {
        return {
           selectedFile: null,
           uploadOutput: []
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
                this.uploadOutput.push(`HTTP ${res.status}: ${res.statusText}`); // Push the response status and statusText
                console.log(res)
            })
            .catch((err) => { 
                console.log(err);
                if (err.response && err.response.status === 400) {
                  this.uploadOutput.push(`HTTP ${err.response.status}: ${err.response.data}`); // Push the error response status and statusText
                  console.log(err.response.status);
                  console.log(err.response.data);
                } else {
                  this.uploadOutput.push(`Error: ${err.message}`); // Push the error message if no response or non-400 status code is available
                }
            })
        }
    }
}
</script>


