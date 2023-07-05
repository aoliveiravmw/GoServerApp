<template>
  <div class="container">
    <div class="logo">
      <img src="../assets/logo.png" alt="Logo">
    </div>
    <div class="server-status">
      <ServerStatus />
    </div>
    <div class="upload-section">
      <UploadFiles />
    </div>
    <div><DisplayPic /></div>
  </div>
</template>

<style>
body {
  background-color: #8f8d8d;
}
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: #f0f6f0;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.logo img {
  width: 150px;
  height: auto;
}

.server-status {
  margin-top: 30px;
  background-color: #a9aba9;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.upload-section {
  margin-top: 30px;
  background-color: #bdc0bd;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.display-pic {
  margin-top: 30px;
  background-color: #c7c9c7;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* Additional Styles for Button */

.button {
  display: inline-block;
  padding: 10px 20px;
  background-color: #4c86c6;
  color: #ffffff;
  border-radius: 4px;
  text-decoration: none;
  transition: background-color 0.3s ease;
}

.button:hover {
  background-color: #3667a9;
}
</style>

<script>
// @ is an alias to /src
import UploadFiles from '@/components/UploadFiles.vue'
import ServerStatus from '../components/ServerStatus.vue'
import DisplayPic from '@/components/DisplayFiles.vue'


export default {
  name: 'SimpleApp',
  components: {
    UploadFiles,
    ServerStatus,
    DisplayPic
  }
}
</script>
