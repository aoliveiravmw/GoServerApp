<template>
  <div id="app">
    <table>
      <tr>
        <th><button @click="downloadFile">Get File</button></th>
        <td>
          <div v-if="fileUrl" class="image-container">
            <img :src="fileUrl" alt="">
          </div>
          <div v-else-if="fileUrl === null">
            <h3>There are no pictures uploaded</h3>
          </div>
          <div v-else>Loading...</div>
        </td>
      </tr>
      <tr>
        <th><button @click="onDelete">Delete File</button></th>
        <td><div>{{ deleteOutput[deleteOutput.length - 1] }}</div></td>
      </tr>
    </table>
  </div>
</template>


<style>
.image-container {
  width: 100%;
  max-width: 500px; /* Adjust the max-width as needed */
  margin: 0 auto;
  text-align: center;
}
.image-container img {
  max-width: 100%;
  height: auto;
}
</style>

<script>
import axios from "axios";

export default {
  name: 'DisplayPic',
  data() {
    return {
      fileUrl: null, // Initialize fileUrl as null
      deleteOutput: []
    }
  },
  methods: {
    downloadFile() {
      axios
        .get('http://192.168.223.23:8081/file', { responseType: 'blob' })
        .then(res => {
          const fileReader = new FileReader();
          fileReader.onload = () => {
            this.fileUrl = fileReader.result;
          };
          fileReader.readAsDataURL(res.data);
        })
        .catch(error => {
          console.error(error);
          this.fileUrl = null; // Set fileUrl to null in case of server error
        });
    },
    onDelete() {
  axios
    .delete('http://192.168.223.23:8081/file')
    .then((res) => {
      this.deleteOutput.push(`HTTP ${res.status}: ${res.statusText}`); // Push the response status and statusText
      console.log(res.status);
      console.log(res.data);
    })
    .catch((err) => { 
      console.log(err);
      if (err.response && err.response.status === 400) {
        this.deleteOutput.push(`HTTP ${err.response.status}: ${err.response.data}`); // Push the error response status and statusText
        console.log(err.response.status);
        console.log(err.response.data);
      } else {
        this.deleteOutput.push(`Error: ${err.message}`); // Push the error message if no response or non-400 status code is available
      }
    });    
}


  },
  mounted() {
    this.downloadFile();
  }
}
</script>
