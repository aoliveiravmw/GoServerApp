<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <UploadFiles/>
  </div>
</template>

<script>
// @ is an alias to /src
import UploadFiles from '@/components/UploadFiles.vue'


export default {
  name: 'AboutView',
  components: {
    UploadFiles
  }
}
</script>
