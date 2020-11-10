import React from 'react';
import { Line as NivoLine } from '@nivo/line';
import addDays from 'date-fns/add_days';
import addHours from 'date-fns/add_hours';
import subDays from 'date-fns/sub_days';
import subHours from 'date-fns/sub_hours';
import dateFormat from 'date-fns/format';
import { useSelector } from 'react-redux';
import PropTypes from 'prop-types';
import './Line.css';

/* todo: check mobile version */
/* todo: fix margin of empty graph */
const Line = ({
    data, color = 'black', width = 275, height = 64,
}) => {
    const interval = useSelector((state) => state.stats.interval);

    return <NivoLine
        enableArea={/* todo: fix error when set to true */ false}
        animate
        enableSlices="x"
        curve="linear"
        debugMesh
        width={width}
        height={height}
        colors={[color]}
        data={data}
        theme={{
            crosshair: {
                line: {
                    /* todo: make solid */
                    stroke: 'black',
                    strokeWidth: 1,
                    strokeOpacity: 0.35,
                },
            },
        }}
        xScale={{
            type: 'linear',
            min: 0,
            max: 'auto',
        }}
        crosshairType={'x'}
        axisLeft={false}
        axisBottom={false}
        enableGridX={false}
        enableGridY={false}
        enablePoints={false}
        xFormat={(index) => {
            if (interval === 1 || interval === 7) {
                const hoursAgo = subHours(Date.now(), 24 * interval);
                return dateFormat(addHours(hoursAgo, index), 'D MMM HH:00');
            }

            const daysAgo = subDays(Date.now(), interval - 1);
            return dateFormat(addDays(daysAgo, index), 'D MMM YYYY');
        }}
        sliceTooltip={(slice) => {
            const { xFormatted, yFormatted } = slice.slice.points[0].data;
            /* todo: check styles */
            return <div className="line__tooltip">
                <span className="line__tooltip-text">
                    <strong>{yFormatted}</strong>
                    <br />
                    <small>{xFormatted}</small>
                </span>
            </div>;
        }}
    />;
};

Line.propTypes = {
    data: PropTypes.array.isRequired,
    color: PropTypes.string,
    width: PropTypes.number,
    height: PropTypes.number,
};

export default Line;
