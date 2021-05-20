<template>
<div>

        <h1 class="h3 mb-3 fw-normal">Enter URL which you want to check</h1>
        <input v-model="Link.Uri" type="text" style="width:70%; height:30px"  placeholder="URL" required autofocus="">
        <input v-model="Link.Level" type="number"  style="width:70px; height:30px; margin:10px" name="num" min=1 max=5 value=1>
        <button class="btn btn-lg btn-primary"   style="height:35px" @click='Check'>Check</button>
        <ul  style="width:50%">
          <li v-for="(errorlink,index) in errorlinks" v-bind:key="errorlink.id">
            {{index}} - {{ errorlink }}
          </li>
        </ul>
        <p>{{NoBr}}</p>
        <span class="progress-text">Loading...</span>

  </div>
</template>

<script>
import axios from 'axios'
export default {
   data: () => ({
   Link: {
        Uri: '',
        Level: 1,
   },
   NoBr:'',
   errorlinks: []
   }),
      methods: {
      Check(){
      this.NoBr=""
      console.log (this.Link)
      this.Link.Level=Number(this.Link.Level)
        axios({url: 'http://localhost:8000/post', data: this.Link, method: 'POST' })
        .then((response) => {

        if (response.data==null)  {
        this.NoBr="No broken Links"
        }
        this.errorlinks = response.data;
          console.log("****", response);
       //   return response;
        })
      }
}
}
</script>

<style>


</style>
