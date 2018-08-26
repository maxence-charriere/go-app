package web

import (
	"github.com/murlokswarm/app"
)

func init() {
	app.Import(&NotFound{})
}

// NotFound is the component that is displayed when the server responds with
// a 404 status.
type NotFound app.ZeroCompo

// Render returns the component markup.
func (n *NotFound) Render() string {
	return `
<div style="
		display: flex;
		flex-direction: column;
    	justify-content: center;
    	align-items: center;
		width: 100%;
		height: 100%;
		overflow: hidden;
		background-color: #21252b;
		color: white;
		font-family: 'Helvetica Neue', 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', sans-serif">
	<h1 style="
		font-size: 100pt; 
		font-weight: 100;
		margin: 0">
		4
		<svg version="1.1" id="Calque_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px"
	 			width="100px" height="100px" viewBox="0 0 186.179 186.18" style="enable-background:new 0 0 186.179 186.18;">
			<circle style="fill:#7FE5D1;" cx="93.089" cy="93.088" r="93.089"/>
			<path style="fill:#79D1C0;" d="M171.303,143.77c-16.638-25.322-45.285-42.047-77.848-42.047c-32.721,0-61.488,16.891-78.088,42.424
				c16.637,25.314,45.285,42.033,77.844,42.033C125.93,186.18,154.701,169.297,171.303,143.77z"/>
			<g>
			<path style="fill:#414244;" d="M84.783,116.143c-1.84,0.834-11.016-5.683-11.928-8.045c-0.907-2.364-3.281-7.36-3.607-10.928
				c-0.322-3.569-4.723-7.08-6.006-11.497c-1.28-4.412,5.799-24.571,7.821-24.571c2.026,0,5.313,15.937,4.885,20.281
				c-0.423,4.348-0.997,11.704-0.741,14.527c0.258,2.827,5.42,7.867,7.007,10.751C83.798,109.546,86.625,115.309,84.783,116.143z"/>
			<path style="fill:#414244;" d="M97.484,115.148c-2.625,0.008-4.008-7.617-4.054-11.114c-0.044-3.499,1.442-9.436,2.761-11.651
				c2.154-3.615,4.512-15.133,7.425-21.256c2.909-6.124,21.12-21.762,22.229-20.649c1.111,1.112-2.168,18.946-6.844,24.537
				c-4.68,5.591-14.064,16.695-15.639,19.236c-1.572,2.538-1.824,9.144-2.123,12.049C100.943,109.207,100.109,115.142,97.484,115.148z
			"/>
			<path style="fill:#414244;" d="M63.714,145.894c-1.749,1.505-7.65-1.265-10.487-4.533c-2.838-3.267-18.779-55.977-19.062-67.453
				C33.882,62.432,44,38.521,45.391,38.521S52.832,57.827,53,67.686c0.172,9.857-4.255,27.355-3.092,34.664
				c1.162,7.306,10.47,25.322,12.073,30.286C63.588,137.603,65.467,144.391,63.714,145.894z"/>
			<path style="fill:#414244;" d="M114.01,129.7c1.346-4.609,14.229-55.01,18.533-60.808c4.309-5.796,30.304-18.072,31.775-16.697
				c1.471,1.375-30.657,43.5-33.475,51.387c-2.822,7.888-6.166,25.391-7.55,29.981c-1.379,4.589-4.985,8.899-7.687,8.138
				C112.904,140.942,112.66,134.308,114.01,129.7z"/>
			</g>
		</svg>
		4
	</h1>
	<p style="
		font-size: 12pt; 
		font-weight: 300;
		margin: 0 0 50px">
		not found
	</p>
</div>
	`
}
