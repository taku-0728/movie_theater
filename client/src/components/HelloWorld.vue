<template>
  <form @submit.prevent="getTheaterList">
    <input type="text" v-model="prefectures" name="prefectures" placeholder="都道府県を入力してください">
    <button type="submit">submit</button>
    <p>{{ message }}</p>
  </form>
</template>

<script>
export default {
  name: 'HelloWorld',
  data() {
    return {
      message: ''
    };
  },  
  methods: {
    async getTheaterList() {
      const result = await this.sendRequest().then((res) => res.text());
      this.message = result;
    },

    async sendRequest() {
      const url = 'http://localhost:8000';
      const data = new URLSearchParams();
      data.append("prefectures",this.prefectures);

      return fetch(url, {
        method: 'POST',
        headers: {
          'X-Requested-With': 'csrf', // csrf header
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: data,
      });
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
