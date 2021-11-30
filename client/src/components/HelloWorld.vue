<template>
  <form @submit.prevent="exec">
    <input type="text" placeholder="都道府県を入力してください">
    <button @click="send">submit</button>
  </form>
</template>

<script>
export default {
  name: 'HelloWorld',
  
  methods: {
    send: function () {
       const url = 'http://localhost:8000';

       fetch(url, {
          method: 'POST',
          headers: {
            'X-Requested-With': 'csrf', // csrf header
            'Content-Type': 'text/html',
          },
          body: 'test',
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
