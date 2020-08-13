export default {
  methods: {
    graphQLErrors(error) {
      return error.graphQLErrors
        .reduce((ax, err) => {
          ax.push(err.message);
          return ax;
        }, [])
        .join(", ");
    },
  },
};
