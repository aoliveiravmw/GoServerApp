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
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.logo {
  margin-bottom: 20px;
}

.server-status {
  margin-bottom: 20px;
}

.upload-section {
  display: flex;
  flex-direction: column;
  align-items: center;
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
