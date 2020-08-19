export default {
  methods: {
    graphQLErrors(error) {
      // eslint-disable-next-line no-console
      if (!error.graphQLErrors) console.log("error", error);
      return error.graphQLErrors
        .reduce((ax, err) => {
          ax.push(err.message);
          return ax;
        }, [])
        .join(", ");
    },
  },
};
