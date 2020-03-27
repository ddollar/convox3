<template>
  <div class="about">
    <h1>This is an about page</h1>
    <h1>Organizations</h1>
    <ul>
      <li v-for="organization in organizations" :key="organization.id">{{ organization.name }}</li>
    </ul>
  </div>
</template>

<script>
import gql from "graphql-tag";

export default {
  apollo: {
    organizations: gql`
      query {
        organizations {
          name
        }
      }
    `
  }
};
</script>