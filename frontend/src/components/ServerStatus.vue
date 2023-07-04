<template>
  <div id="app">
    <table>
      <tr>
        <th>Server status:</th>
        <td>{{ consoleOutput[consoleOutput.length - 1] }}</td>
      </tr>
    </table>
  </div>
</template>



<script>
import axios from "axios";

export default {
  name: 'ServerStatus',
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

table {
  border-collapse: collapse;
  margin: 0 auto;
}

th, td {
  padding: 8px;
  border: none;
}

th {
  font-weight: bold;
  text-align: left;
}

td {
  text-align: center;
}
</style>
