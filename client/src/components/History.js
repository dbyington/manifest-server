import './History.css';

function History(props) {
  const manifests = props.manifests;
  const showManifest = (m) => props.showManifest(m);

  let manifestItems = '';

  if (manifests.length > 0) {
    manifestItems = manifests.map((el, i) => {
      return (
        <li onClick={()=>showManifest(el)} key={i}>{el.title}({el.hashType})</li>
      );
    })
  }

  const manifestList = <ul>{manifestItems}</ul>

  return (
    <div className="History">
      <h3>Cached Manifests</h3>
      {manifestList}
    </div>
  );
}

export default History;
