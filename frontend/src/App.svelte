<script>
  import Router from "svelte-spa-router";
  import Home from "./pages/Home.svelte";
  import NotFound from "./pages/NotFound.svelte";
  import Activate from "./pages/Activate.svelte";
  import Debug from "./pages/Debug.svelte";
    import Account from "./pages/Account.svelte";
    import Settings from "./pages/Settings.svelte";
    import Creator from "./pages/Creator.svelte";
    import Result from "./pages/Result.svelte";

  const routes = {
    // Exact path
    "/": Home,

    "/account": Account,
    "/settings": Settings,
    "/activate": Activate,
    "/creator": Creator,
    "/result": Result,

    "/debug": Debug,

    // Catch-all
    // This is optional, but if present it must be the last
    "*": NotFound,
  };
</script>

<Router {routes} />

