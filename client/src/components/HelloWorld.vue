<template>
  <div>
    <form @submit.prevent="getTheaterList">
      <input type="text" v-model="prefectures" name="prefectures" placeholder="都道府県を入力してください"><br>
      <input type="text" v-model="title" name="title" placeholder="作品名を入力してください"><br>
      <button type="submit">submit</button>
    </form>
    <div v-for='(data, key) in results' :key="key">
      <h3>{{ data.theaterName }}</h3>
      <div v-for='(schedule, key) in data.schedule' :key="key">
        <p>{{ schedule }}</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  data() {
    return {
      results: ''
    };
  },  
  methods: {
    async getTheaterList() {
      const result = await this.sendRequest().then((res) => res.text());
      this.results = JSON.parse(result);
      console.log(this.results);
    },

    async sendRequest() {
      const url = 'http://localhost:8000';
      const data = new URLSearchParams();
      data.append("prefectures",this.prefectures);
      data.append("title",this.title);

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
