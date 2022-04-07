// import logo from './logo.svg';
// import './App.css';

// function App() {
//   return (
//     <div className="App">
//       <header className="App-header">
//         <img src={logo} className="App-logo" alt="logo" />
//         <p>
//           Edit <code>src/App.js</code> and save to reload.
//         </p>
//         <a
//           className="App-link"
//           href="https://reactjs.org"
//           target="_blank"
//           rel="noopener noreferrer"
//         >
//           Learn React
//         </a>
//       </header>
//     </div>
//   );
// }

// export default App;

import React from 'react';
import './App.css';
import Recipe from './Recipe';
class App extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      recipes: []
    }

    this.getRecipes();
  }

  getRecipes() {
    fetch('http://localhost:8080/recipes')
      .then(response => response.json())
      .then(data => this.setState({ recipes: data }));
  }

  render() {
    return (<div>
      {this.state.recipes.map((recipe, index) => (
        <Recipe recipe={recipe} />
      ))}
    </div>);
  }
}
export default App;