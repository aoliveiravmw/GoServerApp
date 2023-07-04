<template>
  <div id="app">
    <table>
      <tr>
        <th>Server status:</th>
        <td>{{ consoleOutput[consoleOutput.length - 1] }}</td>
        <td> 
          <div v-if="StatusBoolean">
            <img class="status-dot" src="../assets/greenDot.png" alt="greenDot">
          </div>
          <div v-else>
            <img class="status-dot" src="../assets/redDot.png" alt="redDot">
          </div>
        </td>
      </tr>
      <tr><button @click="checkServerStatus">Refresh Status</button></tr>
    </table>
  </div>
</template>



<script>
import axios from "axios";

export default {
  name: 'ServerStatus',
  data() {
    return {
      consoleOutput: [],
      StatusBoolean: false 
    }
  },
  methods: {
    checkServerStatus() {
      axios
        .get("http://192.168.223.23:8081/")
        .then((res) => {
          this.consoleOutput.push(res.data);
          this.StatusBoolean = true;
          console.log(res);
        })
        .catch((err) => {
          this.consoleOutput.push(err.message || "Network Error");
          this.StatusBoolean = false;
          console.log(err);
        });
    
    }
  },
  mounted() {
    setInterval(() => this.checkServerStatus(), 30000);
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

th,
td {
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

.status-dot {
  width: 20px; /* Adjust the width as needed */
  height: 20px; /* Adjust the height as needed */
}
</style>