window.addEventListener('load', () => {
  target = document.getElementById('nav_wrapper')

  const openNav = () => {
    target.style.width = "100%";
  }

  const closeNav = () => {
    target.style.width = "0%";
  }

  const toggleNav = query => {
    if (query.matches) {
      openNav()
    } else {
      closeNav()
    }
  }

  document.getElementById('open_nav_button').onclick = openNav
  document.getElementById('close_nav_button').onclick = closeNav
  var watch = window.matchMedia("only screen and (min-width: 768px), (max-aspect-ratio: 5/8)")
  toggleNav(watch)
  watch.addListener(toggleNav)
})