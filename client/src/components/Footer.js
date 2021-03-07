import './Footer.css'

function Footer() {
  const altText = `By clicking 'Submit' on this form you agree that your intention in using this form, the 
  underlying code, systems, and data produced is not malicious. 
  You also agree that if you are found to be using this form, the underlying code, systems, and data produced 
  in any malicious or harmful way you will pay the sum of one thousand dollars to the copywrite owner.`;

  return (
    <div className="Footer">
      <a href={'https://github.com/dbyington/manifest-server'}>Manifest Builder</a> &copy; 2021 Don Byington don!dbyington.com
      <p><small>{altText}</small></p>
    </div>
  );
}

export default Footer;
