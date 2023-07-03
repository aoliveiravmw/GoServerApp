<template>
  <div id="app">
    <img alt="Vue logo" src="./assets/logo.png">
    <h1 class="title"> Server Status </h1>
    <div>{{ consoleOutput[consoleOutput.length - 1] }}</div>
  </div>
</template>


<script>
import axios from "axios";

export default {
  name: 'App',
  data() {
    return {
      consoleOutput: []
    }
  },
  mounted() {
    setInterval(() => {
      axios
        .get("http://192.168.223.23:8081/")
        .then((res) => {
          this.consoleOutput.push(res.data);
          console.log(res.data);
        })
        .catch((err) => {
          this.consoleOutput.push(err.message || "Network Error");
          console.log(err);
        });
    }, 1000);
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
