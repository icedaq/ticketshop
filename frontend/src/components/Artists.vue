<template>
  <div class="hello">
    <h1>Pick an Artist</h1>

    <div class="loading" v-if="loading">
      Loading...
    </div>

    <div v-if="error" class="error">
      {{ error }}
    </div>

      <ul id="artists">
        <li v-for="artist in data">
            <router-link :to="`/shows/${artist.id}`">{{ artist.name }}</router-link>
        </li>
      </ul>
  </div>
</template>

<script>
export default {
  name: 'Artists',
  data () {
    return {
      data: null,
      error: null
    }
  },
  created () {
    // fetch the data when the view is created and the data is
    // already being observed
    this.fetchData()
  },
  watch: {
    // call again the method if the route changes
    '$route': 'fetchData'
  },
  methods: {
    fetchData () {
      this.error = this.post = null
      this.loading = true

      // replace `getPost` with your data fetching util / API wrapper
      this.axios.get(this.$APIServerURL + '/artists').then((response) => {
        this.loading = false
        if (response.status !== 200) {
          this.error = 'Could not load data.'
        } else {
          this.data = response.data
        }
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
  margin: 0;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
  font-size: 200 20px/1.5;
  border-bottom: 1px solid #ccc;
}

li:last-child {
  border: none;
}

a {
  color: #42b983;
}

#artists {
  width: 200px;
  margin: 0 auto;
}

ul {
  list-style-type: none;
  margin: 0;
  padding: 0;
}

li a {
  text-decoration: none;
  color: #000;

  -webkit-transition: font-size 0.3s ease, background-color 0.3s ease;
  -moz-transition: font-size 0.3s ease, background-color 0.3s ease;
  -o-transition: font-size 0.3s ease, background-color 0.3s ease;
  -ms-transition: font-size 0.3s ease, background-color 0.3s ease;
  transition: font-size 0.3s ease, background-color 0.3s ease;
  display: block;
  width: 200px;
}

li a:hover {
  background: #f6f6f6;
}


</style>
